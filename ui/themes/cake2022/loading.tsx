import styles from "./loading.module.scss";

const Loading = () => {
  return (
    <div
      style={{
        height: "100%",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <div
        style={{
          maxWidth: "960px",
          margin: "0 auto",
        }}
      >
        <div className={styles["lds-dual-ring"]}></div>
      </div>
    </div>
  );
};

export default Loading;
