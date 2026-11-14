import { buildSrcset } from '@/api/media'

/**
 * Card sizes attribute: one column on phones, two on tablets, three on
 * desktop — mirrors the HomeSection grid breakpoints.
 */
export const CARD_SIZES = '(max-width: 767px) 100vw, (max-width: 1023px) 50vw, 33vw'

/**
 * Returns the srcset attribute value for an image URL: media-store images
 * get the full width variant set, anything else renders without srcset.
 */
export function imageSrcset(url: string): string {
  return buildSrcset(url)
}
