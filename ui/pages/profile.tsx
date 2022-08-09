import { GetStaticProps } from "next";
import { ProfileUpdateParams } from "props/profile";
import { useEffect, useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import Loading from "../components/loading";
import { api } from "../lib/api";
import useAccount from "../lib/api/account";
import useMessage from "../lib/useMessage";
import ProfileView from "theme/profile";
import { fetchCTF } from "lib/api/ctf";
import { AllPageProps } from "lib/pages";

const Profile = () => {
  const { message, error } = useMessage();
  const { data: account, mutate } = useAccount(null);
  const [country, setCountry] = useState(account?.country || "");
  const { register, setValue, handleSubmit, reset } =
    useForm<ProfileUpdateParams>({
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

  useEffect(() => {
    reset({
      teamname: account?.teamname || "",
    });
    setCountry(account?.country || "");
  }, [account]);

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

export const getStaticProps: GetStaticProps<AllPageProps> = async () => {
  const ctf = await fetchCTF();
  return {
    props: {
      ctf,
    },
  };
};

export default Profile;
