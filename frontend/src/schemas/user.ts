import { z } from 'zod';

// Supondo que a API retorna um usuário com esta estrutura
export const userSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  email: z.string().email(),
  createdAt: z.string().datetime(),
});

export const userListSchema = z.array(userSchema);
