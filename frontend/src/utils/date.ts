/** Renders `dd.mm.yyyy` from an ISO date or a full ISO datetime. */
export function formatDate(iso: string): string {
  const [datePart = iso] = iso.split('T')
  const [year = '', month = '', day = ''] = datePart.split('-')

  return year && month && day ? `${day}.${month}.${year}` : iso
}
