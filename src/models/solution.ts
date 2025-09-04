import z from "zod";

export const SolutionSchema = z.object({
  problem_name: z.string(),
  code: z.string(),
  description: z.string(),
  language: z.string(), 
  problem_link: z.string().optional(),
  problem_id: z.string().optional(),
  difficulty: z.string().optional(),
  notes: z.string().optional(),
});

export type Solution = z.infer<typeof SolutionSchema>;