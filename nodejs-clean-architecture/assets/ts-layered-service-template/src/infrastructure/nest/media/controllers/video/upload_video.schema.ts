import { z } from 'zod'

export const UploadVideoBodySchema = z.object({
  title: z.string().min(1),
})

