import { chakra, FormControl } from "@chakra-ui/react";

const InlineFormControl = chakra(FormControl, {
  baseStyle: {
    display: "flex",
    alignItems: "center",
    label: {
      minWidth: "10em",
    },
  },
});

export default InlineFormControl;
