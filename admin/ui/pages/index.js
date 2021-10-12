import styles from '../styles/Home.module.css'
import Page from "../components/page";

export default function Home() {
    return (
        <Page>
            <div>
                <h1 className={styles.title}>
                    Memo Server Admin Home
                </h1>
            </div>
        </Page>
    )
}
