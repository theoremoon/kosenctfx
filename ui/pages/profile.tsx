import { GetStaticProps } from "next";
import { ProfileUpdateParams } from "props/profile";
import { useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import Loading from "../components/loading";
import { api } from "../lib/api";
import useAccount, { Account, fetchAccount } from "../lib/api/account";
import useMessage from "../lib/useMessage";
import ProfileView from "theme/profile";

interface ProfileProps {
  account: Account | null;
}

const Profile = ({ account: defaultAccount }: ProfileProps) => {
  const { message, error } = useMessage();
  const { data: account, mutate } = useAccount(defaultAccount);
  const [country, setCountry] = useState(account?.country || "");
  const { register, setValue, handleSubmit } = useForm<ProfileUpdateParams>({
    defaultValues: {
      teamname: account?.teamname || "",
      password: "",
    },
  });
  const onSubmit: SubmitHandler<ProfileUpdateParams> = async (data) => {
    try {
      const res = await api.post("/update-profile", {
        ...data,
        country: country,
      });
      message(res);
      mutate();
      setValue("password", "");
    } catch (e) {
      error(e);
    }
  };

  if (account === undefined) {
    return <Loading />;
  }

  return ProfileView({
    register,
    onSubmit: handleSubmit(onSubmit),
    country,
    setCountry,
  });
};

export const getStaticProps: GetStaticProps<ProfileProps> = async () => {
  const account = await fetchAccount().catch(() => null);
  return {
    props: {
      account: account,
    },
  };
};

export default Profile;
