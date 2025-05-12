
### Installing xk6 and Building xk6-kafka
- Requires `Go 1.24.x`
    ```
    go install go.k6.io/xk6/cmd/xk6@latest
    ```
- Build k6 with kafka extension
    ```
    ~/go/bin/xk6 build --with github.com/mostafa/xk6-kafka@latest
    ```
- This will create a k6 binary in current directory with the Kafka extension included.

- Once done, run test using below command
    ```
    ./run.sh
    ```

- Result will be created in `load-test` folder
