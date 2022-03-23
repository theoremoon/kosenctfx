import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  VStack,
} from "@chakra-ui/react";
import CountrySelector from "components/countryselector";
import Right from "components/right";
import { useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import Loading from "../components/loading";
import { api } from "../lib/api";
import useAccount, { Account } from "../lib/api/account";
import useMessage from "../lib/useMessage";

type UpdateParams = {
  teamname: string;
  password: string;
};

interface ProfileProps {
  account: Account;
}

const Profile = ({ account }: ProfileProps) => {
  const { message, error } = useMessage();
  const { mutate } = useAccount();
  const [country, setCountry] = useState(account.country);
  const { register, setValue, handleSubmit } = useForm({
    defaultValues: {
      teamname: account.teamname,
      password: "",
    },
  });
  const onSubmit: SubmitHandler<UpdateParams> = async (data) => {
    try {
      const res = await api.post("/update-profile", {
        ...data,
        country: country,
      });
      message(res);
      mutate();
      setValue("password", "");
    } catch (e) {
      error(e);
    }
  };

  return (
    <Box w="sm" mx="auto" mt="10">
      <form onSubmit={handleSubmit(onSubmit)}>
        <VStack>
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
              {...register("password", { required: false })}
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
              <Button type="submit">Update</Button>
            </Right>
          </FormControl>
        </VStack>
      </form>
    </Box>
  );
};

const ProfileDefault = () => {
  const { data: account } = useAccount();
  if (!account) {
    return <Loading />;
  }
  return <Profile account={account} />;
};

export default ProfileDefault;
