import nekoImage from "./neko.png";
import Image from "next/image";
import { dateFormat } from "lib/date";
import { IndexProps } from "props/index";

const Index = ({ ctf, status }: IndexProps) => {
  return (
    <>
      <section>
        <p>
          {dateFormat(ctf.start_at)} &ndash; {dateFormat(ctf.end_at)}
        </p>
        <p>{status} </p>
      </section>

      <section style={{ display: "flex", justifyContent: "center" }}>
        <Image src={nekoImage} unoptimized={true} />
      </section>

      <h2>[ About ]</h2>
      <section>
        Welcome to CakeCTF 2022! CakeCTF 2022 is a Jeopardy-style Capture The
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
          <li>1st wave: 2022-09-03 14:00:00 JST (UTC+9)</li>
          <li>2nd wave: 2022-09-03 16:00:00 JST (UTC+9)</li>
          <li>Survey: 2022-09-04 02:00:00 JST (UTC+9)</li>
        </ul>
      </section>

      <h2>[ Prize ]</h2>
      <section>
        <p>
          Small gifts will be awarded to teams that
          <ul>
            <li>Solved a specified problem within the second place</li>
            <li>
              Can receive prizes in Japan (Select Japan as your team country to
              be eligible)
            </li>
          </ul>
          <b>One set of prizes</b> will be sent per team. A total five
          challenges, one in each of the five categories (crypto, pwn, web, rev,
          cheat) will be subject to the first/second-blood prizes.
          <b>
            We will announce which challenges are eligible to the
            first/second-blood prizes in{" "}
            <a
              href="https://discord.gg/mP7TqJastk"
              target="_blank"
              rel="noreferrer"
            >
              Discord
            </a>{" "}
          </b>{" "}
          before the CTF starts. All of them will be{" "}
          <b>released in the 2nd wave</b>. (2022-09-03 16:00:00 UTC+9)
          <br />
          We may increase the number of prize slots according to the number of sponsors. <br />
          We cannot send the prize abroad. Please understand at least one of
          your team members needs to reside in Japan to be eligible, as
          mentioned above.
        </p>
        <p>
          以下のチームにはささやかな賞品を用意しています。
          <ul>
            <li>特定の問題を2位以内に解いた</li>
            <li>
              日本国内で賞品を受け取れる（チームの国をJapanに設定してください）
            </li>
          </ul>
          5つのカテゴリ（crypto, pwn, web, rev,
          cheat）の各1問、合計5問がfirst/second-blood賞品の対象です。賞品はチームの人数によらず1問につき1セットです。
          <b>
            どの問題がfirst/second-blood賞品の対象になるかは、{" "}
            <a
              href="https://discord.gg/mP7TqJastk"
              target="_blank"
              rel="noreferrer"
            >
              Discord
            </a>{" "}
            でCTF開始前に公表
          </b>
          します。 これらの問題は<b>2nd waveで公開</b>
          される予定です。（2022-09-03 16:00:00 UTC+9）
          <br />
          スポンサーの数に応じて、賞品の枠は増える可能性があります。<br />
          海外への賞品発送には対応していません。前述した通り、賞品を受け取るにはチームメンバーの少なくとも1人が日本国内に居住している必要があります。
        </p>
      </section>

      <h2>[ Sponsors ]</h2>
      <section>
        <p>
          We&#39;d like to thank the following people for their financial
          support in organizing this event!
        </p>
        <ul>
          <li>atpons</li>
          <li>xrekkusu</li>
          <li>udon</li>
          <li>Edwow Math</li>
          <li>3socha</li>
          <li>kusano_k</li>
          <li>joseph</li>
          <li>y011d4</li>
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
