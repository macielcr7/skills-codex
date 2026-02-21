import { describe, expect, it } from 'vitest'

import { VideoNotFoundError } from '#/domain/repository/media/video/video_not_found.error'

describe('VideoNotFoundError', () => {
  it('sets name and message', () => {
    const err = new VideoNotFoundError('abc')
    expect(err.name).toBe('VideoNotFoundError')
    expect(err.message).toBe('video not found: abc')
  })
})

