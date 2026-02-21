import { z } from 'zod'

import { EntityValidationError } from '#/domain/shared/error/entity_validation_error'

const VideoSchema = z.object({
  id: z.string().uuid(),
  title: z.string().min(1),
})

export type VideoProps = z.infer<typeof VideoSchema>

export class Video {
  public readonly id: string
  public readonly title: string

  constructor(props: VideoProps) {
    this.id = props.id
    this.title = props.title
  }

  validate(): void {
    const result = VideoSchema.safeParse({ id: this.id, title: this.title })
    if (!result.success) {
      throw new EntityValidationError('Video', result.error.issues)
    }
  }
}

