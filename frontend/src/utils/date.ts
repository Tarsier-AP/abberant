export function formatDate(iso: string): string {
  const [year = '', month = '', day = ''] = iso.split('-')

  return year && month && day ? `${day}.${month}.${year}` : iso
}
