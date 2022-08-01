import {
  Box,
  VStack,
  FormControl,
  Input,
  FormLabel,
  Button,
} from "@chakra-ui/react";
import Right from "./components/right";
import CountrySelector from "./components/countrySelector";
import { RegisterProps } from "props/register";

const Register = ({
  register,
  onSubmit,
  country,
  setCountry,
}: RegisterProps) => {
  return (
    <Box w="sm" mx="auto" mt="10">
      <form onSubmit={onSubmit}>
        <VStack>
          <FormControl>
            <FormLabel htmlFor="email">email</FormLabel>
            <Input
              id="email"
              type="email"
              variant="flushed"
              autoComplete="email"
              {...register("email", { required: true })}
            ></Input>
          </FormControl>
          <FormControl>
            <FormLabel htmlFor="teamname">teamname</FormLabel>
            <Input
              id="teamname"
              variant="flushed"
              autoComplete="username"
              {...register("teamname", { required: true })}
            ></Input>
          </FormControl>
          <FormControl>
            <FormLabel htmlFor="password">password</FormLabel>
            <Input
              id="password"
              type="password"
              autoComplete="new-password"
              {...register("password", { required: true })}
            ></Input>
          </FormControl>
          <FormControl>
            <FormLabel htmlFor="country">country</FormLabel>
            <CountrySelector
              id="country"
              value={country}
              onChange={(e) => setCountry(e.target.value)}
            />
          </FormControl>

          <FormControl>
            <Right>
              <Button type="submit">Register</Button>
            </Right>
          </FormControl>
        </VStack>
      </form>
    </Box>
  );
};

export default Register;
