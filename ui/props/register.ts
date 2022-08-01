import { FormEventHandler } from "react";
import { UseFormRegister } from "react-hook-form";

export type RegisterParams = {
  email: string;
  teamname: string;
  password: string;
};

export interface RegisterProps {
  register: UseFormRegister<RegisterParams>;
  onSubmit: FormEventHandler<HTMLFormElement>;
  country: string;
  setCountry: (country: string) => void;
}
