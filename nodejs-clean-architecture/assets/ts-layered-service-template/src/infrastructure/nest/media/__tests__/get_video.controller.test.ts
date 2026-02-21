import { BadRequestException, NotFoundException } from '@nestjs/common'
import { describe, expect, it, vi } from 'vitest'

import { GetVideoController } from '#/infrastructure/nest/media/controllers/video/get_video.controller'
import { VideoNotFoundError } from '#/domain/repository/media/video/video_not_found.error'

describe('GetVideoController', () => {
  it('returns video when usecase succeeds', async () => {
    const getVideo = {
      execute: vi.fn().mockResolvedValue({ id: 'id-1', title: 'hello' }),
    }

    const controller = new GetVideoController(getVideo as never)
    const res = await controller.handle('550e8400-e29b-41d4-a716-446655440000')

    expect(res).toEqual({ id: 'id-1', title: 'hello' })
    expect(getVideo.execute).toHaveBeenCalledWith({
      id: '550e8400-e29b-41d4-a716-446655440000',
    })
  })

  it('throws BadRequestException when id is invalid', async () => {
    const controller = new GetVideoController({ execute: vi.fn() } as never)

    await expect(controller.handle('invalid')).rejects.toBeInstanceOf(
      BadRequestException,
    )
  })

  it('maps VideoNotFoundError to NotFoundException', async () => {
    const getVideo = {
      execute: vi.fn().mockRejectedValue(new VideoNotFoundError('id-1')),
    }

    const controller = new GetVideoController(getVideo as never)

    await expect(
      controller.handle('550e8400-e29b-41d4-a716-446655440000'),
    ).rejects.toBeInstanceOf(NotFoundException)
  })
})
