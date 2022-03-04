import { Box, Flex, Spinner } from "@chakra-ui/react";
import React from "react";

const Loading = () => {
  return (
    <Flex h="100%" justifyContent="center" alignItems="center">
      <Box w="md" mx="auto">
        <Spinner size="xl" />
      </Box>
    </Flex>
  );
};

export default Loading;
