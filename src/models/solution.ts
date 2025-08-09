import z from "zod";

export const SolutionSchema = z.object({
  problemName: z.string(),
  code: z.string(),
  description: z.string(),
  language: z.string(), 
  problemLink: z.string().optional(),
  problemId: z.string().optional(),
  difficulty: z.string().optional(),
  notes: z.string().optional(),
});

export type Solution = z.infer<typeof SolutionSchema>;