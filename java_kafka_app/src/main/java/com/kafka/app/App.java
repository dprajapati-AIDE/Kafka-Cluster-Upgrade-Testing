package com.kafka.app;

public class App {
    public static void main(String[] args) {
        String role = null;
        int msgCount = 10;
        String consumerGroup = "default-group";
        int messageSizeKB = 0;

        for (int i = 0; i < args.length; i++) {
            switch (args[i]) {
                case "--role":
                    if (i + 1 < args.length)
                        role = args[++i];
                    break;
                case "--msg-count":
                    if (i + 1 < args.length)
                        msgCount = Integer.parseInt(args[++i]);
                    break;
                case "--consumer-group":
                    if (i + 1 < args.length)
                        consumerGroup = args[++i];
                    break;
                case "--msg-size-kb":
                    if (i + 1 < args.length)
                        messageSizeKB = Integer.parseInt(args[++i]);
                    break;
                default:
                    System.err.println("Unknown argument: " + args[i]);
            }
        }

        if (role == null) {
            System.err.println("Missing required flag: --role");
            System.exit(1);
        }

        Run.run(role, msgCount, consumerGroup, messageSizeKB);
    }
}
