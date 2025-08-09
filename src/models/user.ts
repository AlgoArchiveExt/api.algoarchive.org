import z from "zod";

export const UserSchema = z.object({
  owner: z.string(), // required
  repoName: z.string(), // required
});

export type User = z.infer<typeof UserSchema>;