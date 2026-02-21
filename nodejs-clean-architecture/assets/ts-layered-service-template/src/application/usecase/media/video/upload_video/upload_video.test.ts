import { describe, expect, it } from 'vitest'

import type { Video } from '#/domain/entity/media/video/video'
import { EntityValidationError } from '#/domain/shared/error/entity_validation_error'
import type { VideoRepository } from '#/domain/repository/media/video/video.repository.interface'

import type { IdGenerator } from '#/application/service/shared/id_generator.interface'
import { UploadVideo } from '#/application/usecase/media/video/upload_video/upload_video'

class FakeVideoRepository implements VideoRepository {
  public saved: Video | null = null

  async save(video: Video): Promise<void> {
    this.saved = video
  }

  async getByID(_id: string): Promise<Video | null> {
    return null
  }
}

class FixedIDGenerator implements IdGenerator {
  constructor(private readonly value: string) {}

  newID(): string {
    return this.value
  }
}

describe('UploadVideo', () => {
  it('creates and saves a Video', async () => {
    const repo = new FakeVideoRepository()
    const idGen = new FixedIDGenerator('550e8400-e29b-41d4-a716-446655440000')

    const usecase = new UploadVideo(repo, idGen)
    const output = await usecase.execute({ title: 'hello' })

    expect(output.id).toBe('550e8400-e29b-41d4-a716-446655440000')
    expect(repo.saved?.id).toBe('550e8400-e29b-41d4-a716-446655440000')
    expect(repo.saved?.title).toBe('hello')
  })

  it('throws when entity validation fails', async () => {
    const repo = new FakeVideoRepository()
    const idGen = new FixedIDGenerator('550e8400-e29b-41d4-a716-446655440000')

    const usecase = new UploadVideo(repo, idGen)
    await expect(usecase.execute({ title: '' })).rejects.toBeInstanceOf(
      EntityValidationError,
    )
  })
})
