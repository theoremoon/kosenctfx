import { Box, chakra } from "@chakra-ui/react";

const Area = chakra(Box, {
  baseStyle: {
    borderWidth: 1,
    borderColor: "white",
    borderRadius: "md",
    p: 4,
    mt: 1,
  },
});

export default Area;
