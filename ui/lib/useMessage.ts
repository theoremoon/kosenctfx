import { useToast } from "@chakra-ui/react";
import axios, { AxiosError, AxiosResponse } from "axios";
import { ErrorResponse, MessageResponse } from "./api";

const useMessage = () => {
  const toast = useToast();
  const err = (e: unknown) => {
    const message = (e as AxiosResponse<ErrorResponse>).data?.message;
    if (message) {
      toast({
        description: message,
        status: "error",
        duration: 2000,
        isClosable: true,
      });
    } else if (axios.isAxiosError(e)) {
      const data = e.response?.data as AxiosError;
      toast({
        description: data.message || JSON.stringify(e.message),
        status: "error",
        duration: 2000,
        isClosable: true,
      });
    }
  };

  const msg = (e: unknown) => {
    const message = (e as AxiosResponse<MessageResponse>).data.message;
    if (message) {
      toast({
        description: message,
        status: "info",
        duration: 2000,
        isClosable: true,
      });
    }
  };

  const text = (m: string) => {
    toast({
      description: m,
      status: "info",
      duration: 2000,
      isClosable: true,
    });
  };

  return { error: err, message: msg, text: text };
};
export default useMessage;
