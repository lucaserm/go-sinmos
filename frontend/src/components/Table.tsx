import type { ReactNode } from 'react'

export function Table({
  head,
  children,
}: {
  head: ReactNode
  children: ReactNode
}) {
  return (
    <div className="overflow-x-auto rounded-xl border border-brand-100 bg-white">
      <table className="w-full text-left text-sm">
        <thead className="border-b border-brand-100 bg-brand-50 text-xs uppercase text-brand-600">
          <tr>{head}</tr>
        </thead>
        <tbody className="divide-y divide-brand-50">{children}</tbody>
      </table>
    </div>
  )
}

export function Th({ children }: { children: ReactNode }) {
  return <th className="px-4 py-3 font-medium">{children}</th>
}

export function Td({ children }: { children: ReactNode }) {
  return <td className="px-4 py-3 text-brand-800">{children}</td>
}
