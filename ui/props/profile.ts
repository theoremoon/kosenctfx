import { Account } from "lib/api/account";
import { FormEventHandler } from "react";
import { UseFormRegister } from "react-hook-form";

export type ProfileUpdateParams = {
  teamname: string;
  password: string;
};

export interface ProfileProps {
  register: UseFormRegister<ProfileUpdateParams>;
  onSubmit: FormEventHandler<HTMLFormElement>;
  country: string;
  setCountry: (country: string) => void;
}
