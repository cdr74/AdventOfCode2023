package cr.aoc2023.day13;

import java.util.Arrays;
import java.util.List;

public class Solver {
    List<Pattern> patterns;
    int part;

    public Solver(List<Pattern> patterns) {
        this.patterns = patterns;
    }

    private int arrayDiff(int[] arr1, int[] arr2) {
        int diff = 0;
        for (int x = 0; x < arr1.length; x++) {
            if (arr1[x] != arr2[x]) diff++;
        }
        return diff;
    }

    private int findMirrorAxis(Pattern pattern) {
        int rows = pattern.getRows();

        for (int x = 0; x + 1 < rows; x++) {
            int delta = x + 1;
            if ((x + delta) >= rows) {
                delta = rows - x - 1;
            }

            int[] before = pattern.getRowMinusRows(x, delta);
            int[] after = pattern.getRowPlusRows(x + 1, delta);

            if (part == 1) {
                if (Arrays.equals(before, after)) return (x + 1);                
            } else {
                if (arrayDiff(before, after) == 1) return (x + 1);
            }
        }

        return 0;
    }

    public int solve() {
        int result = 0;

        for (Pattern pattern : patterns) {
            int mirrorIndex = findMirrorAxis(pattern) * 100;

            if (mirrorIndex == 0) {
                pattern.transposeMatrix();
                mirrorIndex = findMirrorAxis(pattern); // back to original state
                pattern.transposeMatrix();
            }
            result += mirrorIndex;
        }
        return result;
    }
}
