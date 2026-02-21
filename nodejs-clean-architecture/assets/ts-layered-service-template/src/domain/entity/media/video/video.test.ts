import { describe, expect, it } from 'vitest'

import { Video } from '#/domain/entity/media/video/video'
import { EntityValidationError } from '#/domain/shared/error/entity_validation_error'

describe('Video', () => {
  it('validate() passes for valid entity', () => {
    const video = new Video({
      id: '550e8400-e29b-41d4-a716-446655440000',
      title: 'video title',
    })

    expect(() => video.validate()).not.toThrow()
  })

  it('validate() throws for invalid entity', () => {
    const video = new Video({
      // invalid uuid on purpose
      id: 'invalid',
      title: '',
    })

    expect(() => video.validate()).toThrow(EntityValidationError)
  })
})

