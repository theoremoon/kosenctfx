import { PasswordResetRequestProps } from "props/passwordResetRequest";
import Input from "./components/input";
import Button from "./components/button";
import Right from "./components/right";

const ResetRequest = ({ register, onSubmit }: PasswordResetRequestProps) => {
  return (
    <form onSubmit={onSubmit} className="input-form">
      <div className="form-item">
        <label htmlFor="email">email</label>
        <Input
          id="email"
          className="form-input"
          {...register("email", { required: true })}
        />
      </div>
      <div className="form-item">
        <Right>
          <Button type="submit">Send Email</Button>
        </Right>
      </div>
    </form>
  );
};

export default ResetRequest;
