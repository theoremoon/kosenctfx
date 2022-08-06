import axios, { AxiosError, AxiosResponse } from "axios";
import React, { useContext } from "react";
import { ErrorResponse, MessageResponse } from "./api";

export const ToastContext = React.createContext({
  info: (msg: string) => {
    alert(msg);
  },
  error: (msg: string) => {
    alert(msg);
  },
});

const useMessage = () => {
  const { info: infoToast, error: errorToast } = useContext(ToastContext);

  const err = (e: unknown) => {
    const message = (e as AxiosResponse<ErrorResponse>).data?.message;
    if (message) {
      errorToast(message);
    } else if (axios.isAxiosError(e)) {
      const data = e.response?.data as AxiosError;
      errorToast(data.message || JSON.stringify(e.message));
    }
  };

  const msg = (e: unknown) => {
    const message = (e as AxiosResponse<MessageResponse>).data.message;
    if (message) {
      infoToast(message);
    }
  };

  const text = (m: string) => {
    infoToast(m);
  };

  return { error: err, message: msg, text: text };
};
export default useMessage;
