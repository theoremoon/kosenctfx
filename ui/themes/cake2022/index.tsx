import nekoImage from "./neko.png";
import Image from "next/image";
import { dateFormat } from "lib/date";
import { IndexProps } from "props/index";

const Index = ({ ctf, status }: IndexProps) => {
  return (
    <>
      <section>
        <p suppressHydrationWarning={true}>
          {dateFormat(ctf.start_at)} &ndash; {dateFormat(ctf.end_at)}
        </p>
        <p>{status} </p>
      </section>

      <section style={{ display: "flex", justifyContent: "center" }}>
        <Image src={nekoImage} unoptimized={true} alt="" />
      </section>

      <h2>[ About ]</h2>
      <section>
        Welcome to CakeCTF 2023! CakeCTF 2023 is a Jeopardy-style Capture The
        Flag competition hosted by{" "}
        <a
          href="https://twitter.com/y05h1k1ng"
          target="_blank"
          rel="noreferrer"
        >
          yoshiking
        </a>
        ,{" "}
        <a
          href="https://twitter.com/theoremoon"
          target="_blank"
          rel="noreferrer"
        >
          theoremoon
        </a>
        , and{" "}
        <a href="https://twitter.com/ptrYudai" target="_blank" rel="noreferrer">
          ptr-yudai
        </a>
        . There will be challenges in categories such as pwn, web, rev, crypto,
        etc. The challenges are of a difficulty level targeting beginner to
        intermediate players.
        <br />
        This year we have reduced the difficulty level and the number of
        challenges a little more than in previous years. Advanced players are
        encouraged to participate solo or in teams with a couple of people.
      </section>

      <h2>[ Task Release Schedule ]</h2>
      <section>
        <p>
          We will announce the schedule of the challenge release in{" "}
          <a
            href="https://discord.gg/mP7TqJastk"
            target="_blank"
            rel="noreferrer"
          >
            Discord
          </a>
          .
        </p>
        <ul>
          <li>1st wave: 2023-11-11 14:00:00 JST (UTC+9)</li>
          <li>2nd wave: 2023-11-11 16:00:00 JST (UTC+9)</li>
          <li>Survey: 2023-11-12 02:00:00 JST (UTC+9)</li>
        </ul>
      </section>

      <h2>[ Prize ]</h2>
      <section>
        <div>
          <p>
            Small gifts will be awarded to teams that rank in the top 5 who can
            receive prizes in Japan (by selecting Japan as their team country to
            be eligible).
          </p>
          <p>
            We cannot send the prize abroad. Please understand at least one of
            your team members needs to reside in Japan to be eligible, as
            mentioned above.
          </p>
        </div>
        <div>
          <p>
            日本で賞品を受け取れるチームのうち上位5チームにはささやかな賞品が贈られます（チームの国をJapanに設定してください）
          </p>
          <p>
            海外への賞品発送には対応していません。前述した通り、賞品を受け取るにはチームメンバーの少なくとも1人が日本国内に居住している必要があります。
          </p>
        </div>
      </section>

      <h2>[ Sponsors ]</h2>
      <section>
        <p>
          We&#39;d like to thank the following people for their financial
          support in organizing this event!
        </p>
        <ul>
          <li>so</li>
          <li>Edwow Math</li>
          <li>寿司はうまい</li>
          <li>yu1hpa</li>
          <li>udon</li>
          <li>jt</li>
          <li>Xryus Technologies</li>
          <li>沖絢斗</li>
          <li>kusano_k</li>
          <li>kurenaif</li>
        </ul>
        <p>
          Also, the infrastructure of this CTF is sponsored by{" "}
          <a
            href="https://goo.gle/ctfsponsorship"
            target="_blank"
            rel="noreferrer"
          >
            Google CTF Sponsorship
          </a>
          . Thank you!
        </p>
      </section>

      <h2>[ Contact ]</h2>
      <section>
        <p>
          Discord:{" "}
          <a
            href="https://discord.gg/mP7TqJastk"
            target="_blank"
            rel="noreferrer"
          >
            https://discord.gg/mP7TqJastk
          </a>
        </p>
      </section>

      <h2>[ Rules ]</h2>
      <section>
        <ul>
          <li>There is no limit on your team size.</li>
          <li>
            Anyone can participate in this CTF: No restriction on your age,
            nationality, or the editor you use.
          </li>
          <li>
            Your position on the scoreboard is decided by:
            <ol>
              <li>The total points (Higher is better)</li>
              <li>The timestamp of your last submission (Earlier is better)</li>
            </ol>
          </li>
          <li>
            The survey challenge is special: It gives you some points but it
            doesn&#39;t update your &quot;submission timestamp&quot;. You
            can&#39;t get ahead simply by solving the survey faster. Take enough
            time to fill the survey.
          </li>
          <li>
            You can&#39;t brute-force the flag. If you submit 5 incorrect flags
            in a short period of time, the submission form will be locked for 5
            minutes.
          </li>
          <li>You can&#39;t participate in multiple teams.</li>
          <li>
            Sharing the solutions, hints or flags with other teams during the
            competition is strictly forbidden.
          </li>
          <li>You are not allowed to attack the scoreserver.</li>
          <li>You are not allowed to attack the other teams.</li>
          <li>
            You are not allowed to have multiple accounts. If you can&#39;t log
            in to your account, use the password reset form or contact us on
            Discord.
          </li>
          <li>
            We may ban and disqualify any teams that break any of these rules.
          </li>
          <li>
            The flag format is{" "}
            <code>
              CakeCTF\{"{"}[\x20-\x7e]+\{"}"}
            </code>{" "}
            unless specified otherwise.
          </li>
          <li>
            You can ask us in Discord if you have any questions.
            <b>We can't give you hints of the challenges</b>, however.
          </li>
          <li>Most importantly: good luck and have fun!</li>
        </ul>
      </section>
    </>
  );
};
export default Index;
