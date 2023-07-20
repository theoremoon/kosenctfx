import { Box, ChakraProps } from "@chakra-ui/react";
import { lookup } from "country-data";

type CountryFlagProps = {
  country: string;
} & ChakraProps;

const CountryFlag = ({ country, ...props }: CountryFlagProps) => {
  const c = lookup.countries({ alpha2: country })[0];
  if (!c) {
    return <></>;
  }
  return (
    <Box {...props} title={c.name}>
      {c.emoji}
    </Box>
  );
};

export default CountryFlag;
