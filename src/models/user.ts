import z from "zod";

export const UserSchema = z.object({
  owner: z.string(), // required
  repo_name: z.string(), // required
});

export type User = z.infer<typeof UserSchema>;