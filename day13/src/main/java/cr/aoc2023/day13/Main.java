package cr.aoc2023.day13;

import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

public class Main {

    public List<Pattern> readInput(String filename) {
        List<Pattern> patterns = new ArrayList<Pattern>();
        try {
            BufferedReader reader = new BufferedReader(new FileReader(filename));
            String line = reader.readLine();
            List<String> input = new ArrayList<String>();
            while (line != null) {
                if (line.isEmpty()) {
                    patterns.add(new Pattern(input));
                    input.clear();
                } else {
                    input.add(line);
                }
                line = reader.readLine();
            }

            patterns.add(new Pattern(input));
            reader.close();
        } catch (IOException e) {
            e.printStackTrace();
        }
        return patterns;
    }

    public static void main(String[] args) {
        Main main = new Main();
        //List<Pattern> patterns = main.readInput("src/main/resources/test.data");
        List<Pattern> patterns = main.readInput("src/main/resources/actual.data");

        Solver solver = new Solver(patterns);
        solver.part = 1;
        int result1 = solver.solve();
        System.out.println("Part 1: " + result1);

        solver.part = 2;
        int result2 = solver.solve();
        System.out.println("Part 2: " + result2);
    }
}
