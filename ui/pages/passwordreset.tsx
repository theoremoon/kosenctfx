import useMessage from "lib/useMessage";
import { api } from "lib/api";
import { useRouter } from "next/router";
import { ResetParams } from "props/passwordreset";
import { SubmitHandler, useForm } from "react-hook-form";
import PasswordResetView from "theme/passwordreset";

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

export default PasswordReset;
