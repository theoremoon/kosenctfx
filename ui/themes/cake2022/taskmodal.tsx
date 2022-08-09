import { TaskModalProps } from "props/taskmodal";
import Input from "./components/input";
import Button from "./components/button";
import styles from "./taskmodal.module.scss";
import Link from "next/link";

const TaskModal = ({
  task,
  registerFlag,
  onFlagSubmit,
  onClose,
  isSolved,
}: TaskModalProps) => {
  return (
    <>
      <div className={styles["dialog-wrapper"]}>
        <div className={styles["dialog"]}>
          <div className={styles["dialog-left"]}>
            <div className={styles["dialog-head"]}>
              {isSolved && "🎂"}
              {task.name} - {task.score}
            </div>
            <div className={styles["dialog-tags"]}>
              {task.tags.map((tag) => (
                <div className={styles["dialog-tag"]} key={tag}>
                  {tag}
                </div>
              ))}
            </div>
            <div
              className={styles["dialog-description"]}
              dangerouslySetInnerHTML={{ __html: task.description }}
            ></div>
            <div className={styles["dialog-attachments"]}>
              {task.attachments &&
                task.attachments.map((a) => (
                  <a key={a.name} href={a.url} download>
                    {a.name}
                  </a>
                ))}
            </div>
            <div className={styles["dialog-author"]}>author: {task.author}</div>
            <form onSubmit={onFlagSubmit}>
              <Input
                type="text"
                placeholder="CakeCTF{NamuNamu...}"
                {...registerFlag("flag")}
              />
              <Button type="submit">Submit</Button>
            </form>
          </div>
          <div className={styles["dialog-right"]}>
            <div className={styles["dialog-right-inner"]}>
              <div className={styles["dialog-solvenum"]}>
                {task.solved_by.length} solves
              </div>

              <div className={styles["dialog-solves"]}>
                {task.solved_by.map((t) => (
                  <div className={styles["dialog-solve-team"]} key={t.team_id}>
                    <Link href={"/teams/" + t.team_id}>{t.team_name}</Link>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>

      <a onClick={() => onClose()}>
        <div className={styles["dialog-background"]}></div>
      </a>
    </>
  );
};

export default TaskModal;
