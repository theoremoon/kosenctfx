import { chakra, Flex } from "@chakra-ui/react";

const Right = chakra(Flex, {
  baseStyle: {
    flexDirection: "row-reverse",
    w: "100%",
  },
});
export default Right;
