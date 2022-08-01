import { FormEventHandler } from "react";
import { UseFormRegister } from "react-hook-form";

export type LoginParams = {
  teamname: string;
  password: string;
};

export interface LoginProps {
  register: UseFormRegister<LoginParams>;
  onSubmit: FormEventHandler<HTMLFormElement>;
}
