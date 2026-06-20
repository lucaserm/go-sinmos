package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	repo "github.com/lucaserm/go-sinmos/internal/adapters/postgresql/sqlc"
	"github.com/lucaserm/go-sinmos/internal/auth"
	"github.com/lucaserm/go-sinmos/internal/courses"
	"github.com/lucaserm/go-sinmos/internal/enrollments"
	"github.com/lucaserm/go-sinmos/internal/guardians"
	"github.com/lucaserm/go-sinmos/internal/json"
	"github.com/lucaserm/go-sinmos/internal/occurrences"
	occurrencestypes "github.com/lucaserm/go-sinmos/internal/occurrences-types"
	"github.com/lucaserm/go-sinmos/internal/permissions"
	"github.com/lucaserm/go-sinmos/internal/schedules"
	studentguardians "github.com/lucaserm/go-sinmos/internal/student-guardians"
	studentsubjects "github.com/lucaserm/go-sinmos/internal/student-subjects"
	"github.com/lucaserm/go-sinmos/internal/students"
	"github.com/lucaserm/go-sinmos/internal/subjects"
	"github.com/lucaserm/go-sinmos/internal/warnings"
)

type application struct {
	config config
	db     *pgx.Conn
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

func (app *application) mount() *chi.Mux {
	router := chi.NewRouter()

	// middlewares
	router.Use(middleware.RequestID)                              // rate limiting
	router.Use(middleware.ClientIPFromHeader("CF-Connecting-IP")) // rate limiting, analytics and tracing
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	return router
}

func (app *application) v1(r *chi.Mux) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.WriteJSON(w, 200, map[string]string{
			"status": "ok",
		})
	})

	pgRepo := repo.New(app.db)
	authMiddleware := auth.NewMiddleware(pgRepo)

	authService := auth.NewService(pgRepo)
	authHandler := auth.NewHandler(authService)
	authHandler.RegisterRoutes(router, authMiddleware)

	studentService := students.NewService(pgRepo)
	studentHandler := students.NewHandler(studentService)
	studentHandler.RegisterRoutes(router)

	guardianService := guardians.NewService(pgRepo)
	guardianHandler := guardians.NewHandler(guardianService)
	guardianHandler.RegisterRoutes(router)

	courseService := courses.NewService(pgRepo)
	courseHandler := courses.NewHandler(courseService)
	courseHandler.RegisterRoutes(router)

	studentGuardianService := studentguardians.NewService(pgRepo)
	studentGuardianHandler := studentguardians.NewHandler(studentGuardianService)
	studentGuardianHandler.RegisterRoutes(router)

	enrollmentService := enrollments.NewService(pgRepo)
	enrollmentHandler := enrollments.NewHandler(enrollmentService)
	enrollmentHandler.RegisterRoutes(router)

	subjectsService := subjects.NewService(pgRepo)
	subjectsHandler := subjects.NewHandler(subjectsService)
	subjectsHandler.RegisterRoutes(router)

	schedulesService := schedules.NewService(pgRepo)
	schedulesHandler := schedules.NewHandler(schedulesService)
	schedulesHandler.RegisterRoutes(router)

	studentSubjectsService := studentsubjects.NewService(pgRepo)
	studentSubjectsHandler := studentsubjects.NewHandler(studentSubjectsService)
	studentSubjectsHandler.RegisterRoutes(router)

	permissionsService := permissions.NewService(pgRepo)
	permissionsHandler := permissions.NewHandler(permissionsService)
	permissionsHandler.RegisterRoutes(router)

	occurrencesTypesService := occurrencestypes.NewService(pgRepo)
	occurrencesTypesHandler := occurrencestypes.NewHandler(occurrencesTypesService)
	occurrencesTypesHandler.RegisterRoutes(router)

	occurrencesService := occurrences.NewService(pgRepo)
	occurrencesHandler := occurrences.NewHandler(occurrencesService)
	occurrencesHandler.RegisterRoutes(router, authMiddleware)

	warningsService := warnings.NewService(pgRepo)
	warningsHandler := warnings.NewHandler(warningsService)
	warningsHandler.RegisterRoutes(router)

	r.Mount("/api/v1", router)
	return r
}

func (app *application) run(handler http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}
