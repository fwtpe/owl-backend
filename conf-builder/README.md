# MAVEN
https://maven.apache.org/plugins/maven-resources-plugin/

# Build configration
```
../maven/bin/mvn clean process-resources
```
After building  you can find configurations in `target/owl-conf/`.

## Define argument
```
../maven/bin/mvn -D<prop1.name>=<prop value> -D<prop2.name>=<prop value> clean process-resources
```

## Use profile in maven
```
../maven/bin/mvn -s <profile xml> -P <profile id> clean process-resources
```
