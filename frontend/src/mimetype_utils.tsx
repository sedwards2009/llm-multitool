export function isImage(mimetype: string | null): boolean {
  if (mimetype == null) {
    return false;
  }
  return ['image/png', 'image/jpeg', 'image/gif'].includes(mimetype);
}
