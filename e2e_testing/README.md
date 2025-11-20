## How to install

1. Go into this directory

    ```bash
    cd e2e_testing
    ```

2. Create a Python virtual environment

    ```bash
    python -m venv .venv
    ```

3. Activate the virtual environment

    - **Windows (PowerShell):**
        ```bash
        .\.venv\Scripts\Activate
        ```
    - **macOS/Linux:**
        ```bash
        source .venv/bin/activate
        ```

4. Install dependencies

    ```bash
    python -m pip install -r requirements.txt
    ```

5. Run the test
    ```bash
    python test_jobseeker.py
    python test_company_accept_reject.py
    python test_company_post_job.py
    ```
