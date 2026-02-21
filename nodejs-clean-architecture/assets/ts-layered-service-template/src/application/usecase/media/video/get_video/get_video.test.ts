import { describe, expect, it } from 'vitest'

import { Video } from '#/domain/entity/media/video/video'
import type { VideoRepository } from '#/domain/repository/media/video/video.repository.interface'
import { VideoNotFoundError } from '#/domain/repository/media/video/video_not_found.error'

import { GetVideo } from '#/application/usecase/media/video/get_video/get_video'

class FakeVideoRepository implements VideoRepository {
  constructor(private readonly video: Video | null) {}

  async save(_video: Video): Promise<void> {}

  async getByID(_id: string): Promise<Video | null> {
    return this.video
  }
}

describe('GetVideo', () => {
  it('returns video when found', async () => {
    const video = new Video({
      id: '550e8400-e29b-41d4-a716-446655440000',
      title: 'hello',
    })
    video.validate()

    const usecase = new GetVideo(new FakeVideoRepository(video))
    const output = await usecase.execute({ id: video.id })

    expect(output).toEqual({ id: video.id, title: 'hello' })
  })

  it('throws VideoNotFoundError when not found', async () => {
    const usecase = new GetVideo(new FakeVideoRepository(null))
    await expect(
      usecase.execute({ id: '550e8400-e29b-41d4-a716-446655440000' }),
    ).rejects.toBeInstanceOf(VideoNotFoundError)
  })
})
