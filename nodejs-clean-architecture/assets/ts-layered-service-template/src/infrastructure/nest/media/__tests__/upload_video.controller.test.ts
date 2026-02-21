import { BadRequestException } from '@nestjs/common'
import { describe, expect, it, vi } from 'vitest'

import { UploadVideoController } from '#/infrastructure/nest/media/controllers/video/upload_video.controller'
import { EntityValidationError } from '#/domain/shared/error/entity_validation_error'

describe('UploadVideoController', () => {
  it('returns id when usecase succeeds', async () => {
    const uploadVideo = {
      execute: vi.fn().mockResolvedValue({ id: 'id-1' }),
    }

    const controller = new UploadVideoController(uploadVideo as never)
    const res = await controller.handle({ title: 'hello' })

    expect(res).toEqual({ id: 'id-1' })
    expect(uploadVideo.execute).toHaveBeenCalledWith({ title: 'hello' })
  })

  it('throws BadRequestException when body is invalid', async () => {
    const controller = new UploadVideoController({ execute: vi.fn() } as never)

    await expect(controller.handle({ title: '' })).rejects.toBeInstanceOf(
      BadRequestException,
    )
  })

  it('maps EntityValidationError to BadRequestException', async () => {
    const uploadVideo = {
      execute: vi.fn().mockRejectedValue(new EntityValidationError('Video', [])),
    }

    const controller = new UploadVideoController(uploadVideo as never)

    try {
      await controller.handle({ title: 'hello' })
      throw new Error('expected controller to throw')
    } catch (error) {
      expect(error).toBeInstanceOf(BadRequestException)
      const response = (error as BadRequestException).getResponse()
      expect(response).toMatchObject({
        message: 'Video validation failed',
        issues: [],
      })
    }
  })
})
