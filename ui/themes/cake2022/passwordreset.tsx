import { PasswordResetProps } from "props/passwordreset";
import Right from "./components/right";
import Input from "./components/input";
import Button from "./components/button";

const PasswordReset = ({ onSubmit, register }: PasswordResetProps) => {
  return (
    <form onSubmit={onSubmit} className="input-form">
      <div className="form-item">
        <label htmlFor="token">token</label>
        <Input
          id="token"
          className="form-input"
          {...register("token", { required: true })}
        ></Input>
      </div>
      <div className="form-item">
        <label htmlFor="password">password</label>
        <Input
          className="form-input"
          id="password"
          type="password"
          {...register("password", { required: true })}
        ></Input>
      </div>
      <div className="form-item">
        <Right>
          <Button type="submit">Reset Password</Button>
        </Right>
      </div>
    </form>
  );
};

export default PasswordReset;
