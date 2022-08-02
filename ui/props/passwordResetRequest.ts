import { FormEventHandler } from "react";
import { UseFormRegister } from "react-hook-form";

export type ResetRequestParams = {
  email: string;
};

export interface PasswordResetRequestProps {
  register: UseFormRegister<ResetRequestParams>;
  onSubmit: FormEventHandler<HTMLFormElement>;
}
