import { z } from 'zod'

export const GetVideoParamsSchema = z.object({
  id: z.string().uuid(),
})

