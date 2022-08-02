import { api } from "lib/api";
import { fetchAccount } from "lib/api/account";
import { fetchCTF } from "lib/api/ctf";
import { AllPageProps } from "lib/pages";
import { isStaticMode } from "lib/static";
import useMessage from "lib/useMessage";
import { GetStaticProps } from "next";
import { useRouter } from "next/router";
import { ResetRequestParams } from "props/passwordResetRequest";
import { SubmitHandler, useForm } from "react-hook-form";
import ResetRequestView from "theme/passwordResetRequest";

const ResetRequest = () => {
  const { register, handleSubmit } = useForm<ResetRequestParams>();
  const router = useRouter();
  const { error: errorMessage } = useMessage();
  const onSubmit: SubmitHandler<ResetRequestParams> = async (values) => {
    try {
      await api.post("/passwordreset-request", {
        email: values.email,
      });
      router.push("/passwordreset");
    } catch (e) {
      errorMessage(e);
    }
  };
  return ResetRequestView({
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

export default ResetRequest;
