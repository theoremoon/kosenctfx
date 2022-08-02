import { FormEventHandler } from "react";
import { UseFormRegister } from "react-hook-form";

export type ResetParams = {
  token: string;
  password: string;
};
export interface PasswordResetProps {
  register: UseFormRegister<ResetParams>;
  onSubmit: FormEventHandler<HTMLFormElement>;
}
