import useMessage from "lib/useMessage";
import { api } from "lib/api";
import { useRouter } from "next/router";
import { ResetParams } from "props/passwordreset";
import { SubmitHandler, useForm } from "react-hook-form";
import PasswordResetView from "theme/passwordreset";
import { isStaticMode } from "lib/static";
import { fetchCTF } from "lib/api/ctf";
import { GetStaticProps } from "next";
import { AllPageProps } from "lib/pages";
import { fetchAccount } from "lib/api/account";

const PasswordReset = () => {
  const { register, handleSubmit } = useForm<ResetParams>();
  const router = useRouter();
  const { error: errorMessage } = useMessage();
  const onSubmit: SubmitHandler<ResetParams> = async (values) => {
    try {
      await api.post("/passwordreset", {
        token: values.token,
        new_password: values.password,
      });

      router.push("/login");
    } catch (e) {
      errorMessage(e);
    }
  };
  return PasswordResetView({
    register,
    onSubmit: handleSubmit(onSubmit),
  });
};

export const getStaticProps: GetStaticProps<AllPageProps> = async () => {
  const account = isStaticMode ? null : await fetchAccount().catch(() => null);
  const ctf = await fetchCTF();
  return {
    props: {
      account: account,
      ctf: ctf,
    },
  };
};

export default PasswordReset;
