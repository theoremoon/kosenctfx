import FormWrapper from "../../components/formwrapper";
import FormItem from "../../components/formitem";
import InlineLabel from "../../components/inlinelabel";
import InlineInput from "../../components/inlineinput";
import InlineCheck from "../../components/inlinecheck";
import Right from "../../components/right";
import Button from "../../components/button";
import {useForm} from "react-hook-form";
import {api} from "lib/api";
import {useEffect} from "react";

interface IBucketConfig {
  endpoint: string
  region: string|null
  bucketName: string
  accessKey: string
  secretKey: string
  https: boolean
}

const Config = () => {
  const { register, handleSubmit, setValue } = useForm<IBucketConfig>();
  const onSubmit = async (data: IBucketConfig) => {
    try {
      const res = await api.post("/admin/set-bucket", data);
      console.log(res);
    }catch(e) {
      console.log(e);
    }
  }

  useEffect(() => {
    api.get<IBucketConfig>("/admin/get-bucket").then(res => {
      setValue("endpoint", res.data.endpoint);
      setValue("region", res.data.region);
      setValue("bucketName", res.data.bucketName);
      setValue("accessKey", res.data.accessKey);
      setValue("secretKey", res.data.secretKey);
      setValue("https", res.data.https);
    }).catch(e => {
      console.log(e);
    })
  }, [])

  return (
    <FormWrapper>
      <form onSubmit={handleSubmit(onSubmit)}>
        <FormItem>
          <InlineLabel>Endpoint</InlineLabel>
          <InlineInput type="text" placeholder="s3.amazon.com" {...register("endpoint", {
            required: true,
          })} />

          <InlineLabel>Region</InlineLabel>
          <InlineInput type="text" placeholder="ap-northeast-1" {...register("region", {
            required: false,
          })}  />

          <InlineLabel>Bucket Name</InlineLabel>
          <InlineInput type="text" placeholder=""  {...register("bucketName", {
            required: true,
          })} />

          <InlineLabel>Access Key</InlineLabel>
          <InlineInput type="text"  {...register("accessKey", {
            required: true,
          })} />

          <InlineLabel>Secret Key</InlineLabel>
          <InlineInput type="password"  {...register("secretKey", {
            required: true,
          })} />

          <InlineLabel>HTTPS</InlineLabel>
          <InlineCheck selected={true}  {...register("https", {
          })} />
        </FormItem>

        <FormItem>
          <Right>
            <Button>apply</Button>
          </Right>
        </FormItem>
      </form>
    </FormWrapper>
  );
};

export default Config;
