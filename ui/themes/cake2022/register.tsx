import CountrySelector from "./components/countrySelector";
import { RegisterProps } from "props/register";
import Right from "./components/right";
import Input from "./components/input";
import Button from "./components/button";

const Register = ({
  register,
  onSubmit,
  country,
  setCountry,
}: RegisterProps) => {
  return (
    <form onSubmit={onSubmit} className="input-form">
      <div className="form-item">
        <label htmlFor="email">email</label>
        <Input
          id="email"
          className="form-input"
          type="email"
          autoComplete="email"
          {...register("email", { required: true })}
        ></Input>
      </div>
      <div className="form-item">
        <label htmlFor="teamname">teamname</label>
        <Input
          id="teamname"
          className="form-input"
          autoComplete="username"
          {...register("teamname", { required: true })}
        ></Input>
      </div>
      <div className="form-item">
        <label htmlFor="password">password</label>
        <Input
          id="password"
          type="password"
          className="form-input"
          autoComplete="new-password"
          {...register("password", { required: true })}
        ></Input>
      </div>
      <div className="form-item">
        <label htmlFor="country">country</label>
        <CountrySelector
          id="country"
          className="form-input"
          value={country}
          onChange={(e) => setCountry(e.target.value)}
        />
      </div>

      <div className="form-item">
        <Right>
          <Button type="submit">Register</Button>
        </Right>
      </div>
    </form>
  );
};

export default Register;
