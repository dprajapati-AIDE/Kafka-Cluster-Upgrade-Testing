# Setting up Multiple Java Version

## Installation and Configuration of jenv

Install `jenv`
```
brew install jenv
```

Add the following lines to your shell startup file
```
export PATH="$HOME/.jenv/bin:$PATH"
eval "$(jenv init -)"
```

Then, source the file to apply the changes:
```
source ~/.zshrc
```


## Installation of Java Versions

Install Java8 using below command
```
brew install --cask temurin@8
```

Install Java11 using below command
```
brew install openjdk@11
```


## Configuration

Add Java8 to jenv
```
jenv add /Library/Java/JavaVirtualMachines/temurin-8.jdk/Contents/Home
```

Setup symlink for Java11 and add to jenv
```
sudo ln -sfn /opt/homebrew/opt/openjdk@11/libexec/openjdk.jdk /Library/Java/JavaVirtualMachines/openjdk-11.jdk
```

```
jenv add /opt/homebrew/Cellar/openjdk@11/11.0.27/libexec/openjdk.jdk/Contents/Home
```

#### List all the versions
```
jenv versions
```

#### List current version
```
jenv version
```

## Set Global and Local versions
The global version will be the default for your system, while the local version will be specific to the current directory.

```
jenv global 11.0
```

```
jenv local 11.0
```

## Create Maven Project

Update details like groupId and artifactId based on requirement
```
mvn archetype:generate -DgroupId=com.kafka.app -DartifactId=java_kafka_app -DarchetypeArtifactId=maven-archetype-quickstart -DinteractiveMode=false
```

### To compile for Java 8 (default):
```
jenv local 1.8
```
```
mvn clean install
```

### To compile for Java 11:
```
jenv local 11
```
```
mvn clean install -Djava.version=11
```

### To compile for Java 21:
```
jenv local 21
```
```
mvn clean install -Djava.version=21
```
---

#### VS Code Classpath issue
- Go to settings or open it using `command + ,`
- Search for `java.project.sourcePaths`
- Add `"src/main/java"`
- Reload the workspace.
