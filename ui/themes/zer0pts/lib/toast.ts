import { IToast, useToast } from "@chakra-ui/react";

type ToastF = (option: IToast) => void;
const toast = (toast: ToastF) => {
  return {
    info: (msg: string) => {
      toast({
        description: msg,
        status: "info",
        duration: 2000,
        isClosable: true,
      });
    },
    error: (msg: string) => {
      toast({
        description: msg,
        status: "error",
        duration: 2000,
        isClosable: true,
      });
    },
  };
};

export default toast;
